import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:rupu/domain/entities/voting.dart';
import 'package:rupu/domain/entities/voting_results_result.dart';
import 'package:rupu/presentation/views/horizontal_properties/controllers/voting_group_detail_controller.dart';
import 'package:rupu/presentation/widgets/shared/card_actions.dart';
import 'package:go_router/go_router.dart';

class VotingGroupDetailView
    extends GetWidget<VotingGroupDetailController> {
  final int propertyId;
  final int votingGroupId;

  const VotingGroupDetailView({
    super.key,
    required this.propertyId,
    required this.votingGroupId,
  });

  @override
  String? get tag =>
      VotingGroupDetailController.tagFor(propertyId, votingGroupId);

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Obx(() => Text(controller.group.value?.name ?? 'Detalle de grupo')),
      ),
      body: Obx(
        () {
          if (controller.isLoading.value) {
            return const Center(child: CircularProgressIndicator());
          }
          if (controller.errorMessage.value != null) {
            return Center(child: Text(controller.errorMessage.value!));
          }
          return ListView.builder(
            itemCount: controller.votings.length,
            itemBuilder: (context, index) {
              final voting = controller.votings[index];
              return _VotingCard(voting: voting);
            },
          );
        },
      ),
    );
  }
}

class _VotingFormBottomSheet extends StatefulWidget {
  final String title;
  final String actionLabel;
  final Voting? voting;
  final Future<bool> Function(
    Map<String, dynamic> data,
  ) onSubmit;

  const _VotingFormBottomSheet({
    required this.title,
    required this.actionLabel,
    required this.onSubmit,
    this.voting,
  });

  @override
  State<_VotingFormBottomSheet> createState() => _VotingFormBottomSheetState();
}

class _VotingFormBottomSheetState extends State<_VotingFormBottomSheet> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _titleCtrl;
  late final TextEditingController _descriptionCtrl;
  bool _isSecret = false;
  bool _allowAbstention = true;
  bool _saving = false;
  String? _errorMessage;

  @override
  void initState() {
    super.initState();
    final voting = widget.voting;

    _titleCtrl = TextEditingController(text: voting?.title ?? '');
    _descriptionCtrl =
        TextEditingController(text: voting?.description ?? '');
    _isSecret = voting?.isSecret ?? false;
    _allowAbstention = voting?.allowAbstention ?? true;
  }

  @override
  void dispose() {
    _titleCtrl.dispose();
    _descriptionCtrl.dispose();
    super.dispose();
  }

  Map<String, dynamic> _buildPayload() {
    final payload = <String, dynamic>{
      'title': _titleCtrl.text.trim(),
      'description': _descriptionCtrl.text.trim(),
      'is_secret': _isSecret,
      'allow_abstention': _allowAbstention,
    };

    payload.removeWhere((key, value) => value == null);
    return payload;
  }

  Future<void> _handleSubmit() async {
    if (_saving) return;
    final formState = _formKey.currentState;
    if (formState == null || !formState.validate()) {
      return;
    }

    FocusScope.of(context).unfocus();
    setState(() {
      _saving = true;
      _errorMessage = null;
    });

    try {
      final result = await widget.onSubmit(_buildPayload());
      if (!mounted) return;
      if (!result) {
        setState(() {
          _saving = false;
          _errorMessage = 'No se pudo guardar la votación. Inténtalo más tarde.';
        });
        return;
      }
      Navigator.of(context).pop(result);
    } catch (_) {
      if (!mounted) return;
      setState(() {
        _saving = false;
        _errorMessage =
            'Ocurrió un error al guardar la votación. Inténtalo nuevamente.';
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final tt = Theme.of(context).textTheme;
    final viewInsets = MediaQuery.viewInsetsOf(context);

    InputDecoration decoration(String label, {String? hint}) {
      return InputDecoration(
        labelText: label,
        hintText: hint,
        filled: true,
        fillColor: cs.surfaceContainerHighest,
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(14)),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: BorderSide(color: cs.outlineVariant),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: BorderSide(color: cs.primary.withValues(alpha: .5)),
        ),
      );
    }

    Widget field({
      required String label,
      required TextEditingController controller,
      String? hint,
      TextInputType keyboardType = TextInputType.text,
      String? Function(String?)? validator,
    }) {
      return TextFormField(
        controller: controller,
        decoration: decoration(label, hint: hint),
        keyboardType: keyboardType,
        textInputAction: TextInputAction.next,
        validator: validator,
      );
    }

    return FractionallySizedBox(
      heightFactor: 0.95,
      child: SafeArea(
        top: false,
        child: Padding(
          padding: EdgeInsets.fromLTRB(16, 0, 16, viewInsets.bottom + 16),
          child: DecoratedBox(
            decoration: BoxDecoration(
              color: cs.surface,
              borderRadius: const BorderRadius.only(
                topLeft: Radius.circular(32),
                topRight: Radius.circular(32),
              ),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withValues(alpha: .18),
                  blurRadius: 28,
                  offset: const Offset(0, 22),
                ),
              ],
            ),
            child: ClipRRect(
              borderRadius: const BorderRadius.only(
                topLeft: Radius.circular(32),
                topRight: Radius.circular(32),
              ),
              child: Column(
                children: [
                  const SizedBox(height: 12),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(24, 10, 24, 20),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          widget.title,
                          style: tt.titleMedium?.copyWith(
                            fontWeight: FontWeight.w800,
                          ),
                        ),
                        const SizedBox(height: 6),
                        Text(
                          'Completa la información de la votación para continuar.',
                          style: tt.bodyMedium?.copyWith(
                            color: cs.onSurfaceVariant,
                          ),
                        ),
                      ],
                    ),
                  ),
                  Expanded(
                    child: SingleChildScrollView(
                      padding: const EdgeInsets.fromLTRB(24, 0, 24, 24),
                      child: Form(
                        key: _formKey,
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            field(
                              label: 'Título',
                              controller: _titleCtrl,
                              validator: (value) {
                                if (value == null || value.trim().isEmpty) {
                                  return 'Ingresa el título de la votación';
                                }
                                return null;
                              },
                            ),
                            const SizedBox(height: 16),
                            field(
                              label: 'Descripción',
                              controller: _descriptionCtrl,
                            ),
                            const SizedBox(height: 8),
                            SwitchListTile.adaptive(
                              value: _isSecret,
                              onChanged: _saving
                                  ? null
                                  : (value) {
                                      setState(() {
                                        _isSecret = value;
                                      });
                                    },
                              title: const Text('Votación secreta'),
                            ),
                            SwitchListTile.adaptive(
                              value: _allowAbstention,
                              onChanged: _saving
                                  ? null
                                  : (value) {
                                      setState(() {
                                        _allowAbstention = value;
                                      });
                                    },
                              title: const Text('Permitir abstención'),
                            ),
                            if (_errorMessage != null) ...[
                              const SizedBox(height: 16),
                              Container(
                                width: double.infinity,
                                padding: const EdgeInsets.all(12),
                                decoration: BoxDecoration(
                                  color: cs.errorContainer,
                                  borderRadius: BorderRadius.circular(14),
                                ),
                                child: Row(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Icon(Icons.error_outline, color: cs.error),
                                    const SizedBox(width: 12),
                                    Expanded(
                                      child: Text(
                                        _errorMessage!,
                                        style: tt.bodyMedium?.copyWith(
                                          color: cs.onErrorContainer,
                                          fontWeight: FontWeight.w600,
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ],
                            const SizedBox(height: 24),
                            Row(
                              children: [
                                Expanded(
                                  child: OutlinedButton(
                                    onPressed: _saving
                                        ? null
                                        : () => Navigator.of(context).pop(),
                                    child: const Text('Cancelar'),
                                  ),
                                ),
                                const SizedBox(width: 12),
                                Expanded(
                                  child: FilledButton.icon(
                                    onPressed: _saving ? null : _handleSubmit,
                                    icon: _saving
                                        ? SizedBox(
                                            width: 18,
                                            height: 18,
                                            child: CircularProgressIndicator(
                                              strokeWidth: 2.2,
                                              color: cs.onPrimary,
                                            ),
                                          )
                                        : const Icon(Icons.save_outlined),
                                    label: Text(
                                      _saving
                                          ? 'Guardando...'
                                          : widget.actionLabel,
                                    ),
                                  ),
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _VotingCard extends StatefulWidget {
  final Voting voting;

  const _VotingCard({required this.voting});

  @override
  State<_VotingCard> createState() => _VotingCardState();
}

class _VotingCardState extends State<_VotingCard> {
  bool _isExpanded = false;

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final tt = Theme.of(context).textTheme;
    final controller = Get.find<VotingGroupDetailController>();

    return Card(
      margin: const EdgeInsets.all(8.0),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(widget.voting.title, style: tt.titleMedium),
            const SizedBox(height: 8.0),
            Text(widget.voting.description, style: tt.bodyMedium),
            const SizedBox(height: 8.0),
            Row(
              children: [
                Text(widget.voting.isActive ? 'Activa' : 'Inactiva'),
                const Spacer(),
                IconButton(
                  icon: Icon(
                      _isExpanded ? Icons.expand_less : Icons.expand_more),
                  onPressed: () {
                    setState(() {
                      _isExpanded = !_isExpanded;
                    });
                    if (_isExpanded) {
                      controller.getVotingOptions(widget.voting.id);
                    }
                  },
                ),
                PopupMenuButton<String>(
                  onSelected: (value) {
                    if (value == 'live') {
                      final state = GoRouter.of(context)
                          .routerDelegate
                          .currentConfiguration;
                      final segments = state.uri.pathSegments;
                      final page = segments.length > 1 ? segments[1] : '0';
                      final propertyId =
                          Get.find<VotingGroupDetailController>().propertyId;
                      final path =
                          '/home/$page/horizontal-properties/$propertyId/voting/${widget.voting.votingGroupId}/votings/${widget.voting.id}/live';
                      GoRouter.of(context).push(path);
                    } else if (value == 'edit') {
                      showModalBottomSheet<bool>(
                        context: context,
                        backgroundColor: Colors.transparent,
                        useRootNavigator: true,
                        isScrollControlled: true,
                        builder: (_) => _VotingFormBottomSheet(
                          title: 'Editar votación',
                          actionLabel: 'Guardar cambios',
                          voting: widget.voting,
                          onSubmit: (payload) =>
                              Get.find<VotingGroupDetailController>()
                                  .updateVoting(
                            votingId: widget.voting.id,
                            data: payload,
                          ),
                        ),
                      );
                    } else if (value == 'delete') {
                      Get.find<VotingGroupDetailController>()
                          .deleteVoting(widget.voting.id);
                    } else if (value == 'toggle') {
                      if (widget.voting.isActive) {
                        Get.find<VotingGroupDetailController>()
                            .deactivateVoting(widget.voting.id);
                      } else {
                        Get.find<VotingGroupDetailController>()
                            .activateVoting(widget.voting.id);
                      }
                    }
                  },
                  itemBuilder: (BuildContext context) =>
                      <PopupMenuEntry<String>>[
                    const PopupMenuItem<String>(
                      value: 'live',
                      child: Text('Votación en vivo'),
                    ),
                    const PopupMenuItem<String>(
                      value: 'edit',
                      child: Text('Editar'),
                    ),
                    const PopupMenuItem<String>(
                      value: 'delete',
                      child: Text('Eliminar'),
                    ),
                    PopupMenuItem<String>(
                      value: 'toggle',
                      child: Text(
                          widget.voting.isActive ? 'Desactivar' : 'Activar'),
                    ),
                  ],
                ),
              ],
            ),
            if (_isExpanded) _VotingOptions(voting: widget.voting),
          ],
        ),
      ),
    );
  }
}

class _VotingOptionFormBottomSheet extends StatefulWidget {
  final String title;
  final String actionLabel;
  final Future<bool> Function(
    Map<String, dynamic> data,
  ) onSubmit;

  const _VotingOptionFormBottomSheet({
    required this.title,
    required this.actionLabel,
    required this.onSubmit,
  });

  @override
  State<_VotingOptionFormBottomSheet> createState() =>
      _VotingOptionFormBottomSheetState();
}

class _VotingOptionFormBottomSheetState
    extends State<_VotingOptionFormBottomSheet> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _optionTextCtrl;
  late final TextEditingController _optionCodeCtrl;
  late final TextEditingController _colorCtrl;
  bool _saving = false;
  String? _errorMessage;

  @override
  void initState() {
    super.initState();
    _optionTextCtrl = TextEditingController();
    _optionCodeCtrl = TextEditingController();
    _colorCtrl = TextEditingController();
  }

  @override
  void dispose() {
    _optionTextCtrl.dispose();
    _optionCodeCtrl.dispose();
    _colorCtrl.dispose();
    super.dispose();
  }

  Map<String, dynamic> _buildPayload() {
    final payload = <String, dynamic>{
      'option_text': _optionTextCtrl.text.trim(),
      'option_code': _optionCodeCtrl.text.trim(),
      'color': _colorCtrl.text.trim(),
    };

    payload.removeWhere((key, value) => value == null);
    return payload;
  }

  Future<void> _handleSubmit() async {
    if (_saving) return;
    final formState = _formKey.currentState;
    if (formState == null || !formState.validate()) {
      return;
    }

    FocusScope.of(context).unfocus();
    setState(() {
      _saving = true;
      _errorMessage = null;
    });

    try {
      final result = await widget.onSubmit(_buildPayload());
      if (!mounted) return;
      if (!result) {
        setState(() {
          _saving = false;
          _errorMessage =
              'No se pudo guardar la opción de votación. Inténtalo más tarde.';
        });
        return;
      }
      Navigator.of(context).pop(result);
    } catch (_) {
      if (!mounted) return;
      setState(() {
        _saving = false;
        _errorMessage =
            'Ocurrió un error al guardar la opción de votación. Inténtalo nuevamente.';
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final cs = Theme.of(context).colorScheme;
    final tt = Theme.of(context).textTheme;
    final viewInsets = MediaQuery.viewInsetsOf(context);

    InputDecoration decoration(String label, {String? hint}) {
      return InputDecoration(
        labelText: label,
        hintText: hint,
        filled: true,
        fillColor: cs.surfaceContainerHighest,
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(14)),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: BorderSide(color: cs.outlineVariant),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(14),
          borderSide: BorderSide(color: cs.primary.withValues(alpha: .5)),
        ),
      );
    }

    Widget field({
      required String label,
      required TextEditingController controller,
      String? hint,
      TextInputType keyboardType = TextInputType.text,
      String? Function(String?)? validator,
    }) {
      return TextFormField(
        controller: controller,
        decoration: decoration(label, hint: hint),
        keyboardType: keyboardType,
        textInputAction: TextInputAction.next,
        validator: validator,
      );
    }

    return FractionallySizedBox(
      heightFactor: 0.95,
      child: SafeArea(
        top: false,
        child: Padding(
          padding: EdgeInsets.fromLTRB(16, 0, 16, viewInsets.bottom + 16),
          child: DecoratedBox(
            decoration: BoxDecoration(
              color: cs.surface,
              borderRadius: const BorderRadius.only(
                topLeft: Radius.circular(32),
                topRight: Radius.circular(32),
              ),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withValues(alpha: .18),
                  blurRadius: 28,
                  offset: const Offset(0, 22),
                ),
              ],
            ),
            child: ClipRRect(
              borderRadius: const BorderRadius.only(
                topLeft: Radius.circular(32),
                topRight: Radius.circular(32),
              ),
              child: Column(
                children: [
                  const SizedBox(height: 12),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(24, 10, 24, 20),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          widget.title,
                          style: tt.titleMedium?.copyWith(
                            fontWeight: FontWeight.w800,
                          ),
                        ),
                        const SizedBox(height: 6),
                        Text(
                          'Completa la información de la opción para continuar.',
                          style: tt.bodyMedium?.copyWith(
                            color: cs.onSurfaceVariant,
                          ),
                        ),
                      ],
                    ),
                  ),
                  Expanded(
                    child: SingleChildScrollView(
                      padding: const EdgeInsets.fromLTRB(24, 0, 24, 24),
                      child: Form(
                        key: _formKey,
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            field(
                              label: 'Texto de la opción',
                              controller: _optionTextCtrl,
                              validator: (value) {
                                if (value == null || value.trim().isEmpty) {
                                  return 'Ingresa el texto de la opción';
                                }
                                return null;
                              },
                            ),
                            const SizedBox(height: 16),
                            field(
                              label: 'Código de la opción',
                              controller: _optionCodeCtrl,
                            ),
                            const SizedBox(height: 16),
                            field(
                              label: 'Color',
                              controller: _colorCtrl,
                            ),
                            if (_errorMessage != null) ...[
                              const SizedBox(height: 16),
                              Container(
                                width: double.infinity,
                                padding: const EdgeInsets.all(12),
                                decoration: BoxDecoration(
                                  color: cs.errorContainer,
                                  borderRadius: BorderRadius.circular(14),
                                ),
                                child: Row(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    Icon(Icons.error_outline, color: cs.error),
                                    const SizedBox(width: 12),
                                    Expanded(
                                      child: Text(
                                        _errorMessage!,
                                        style: tt.bodyMedium?.copyWith(
                                          color: cs.onErrorContainer,
                                          fontWeight: FontWeight.w600,
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ],
                            const SizedBox(height: 24),
                            Row(
                              children: [
                                Expanded(
                                  child: OutlinedButton(
                                    onPressed: _saving
                                        ? null
                                        : () => Navigator.of(context).pop(),
                                    child: const Text('Cancelar'),
                                  ),
                                ),
                                const SizedBox(width: 12),
                                Expanded(
                                  child: FilledButton.icon(
                                    onPressed: _saving ? null : _handleSubmit,
                                    icon: _saving
                                        ? SizedBox(
                                            width: 18,
                                            height: 18,
                                            child: CircularProgressIndicator(
                                              strokeWidth: 2.2,
                                              color: cs.onPrimary,
                                            ),
                                          )
                                        : const Icon(Icons.save_outlined),
                                    label: Text(
                                      _saving
                                          ? 'Guardando...'
                                          : widget.actionLabel,
                                    ),
                                  ),
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _VotingOptions extends StatelessWidget {
  final Voting voting;

  const _VotingOptions({required this.voting});

  @override
  Widget build(BuildContext context) {
    final controller = Get.find<VotingGroupDetailController>();
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Divider(),
        const Text('Opciones de votación'),
        Obx(
          () => Column(
            children: controller.options
                .map(
                  (option) => ListTile(
                    title: Text(option.optionText),
                    trailing: IconButton(
                      icon: const Icon(Icons.delete),
                      onPressed: () {
                        controller.deleteVotingOption(
                          votingId: voting.id,
                          optionId: option.id,
                        );
                      },
                    ),
                  ),
                )
                .toList(),
          ),
        ),
        ElevatedButton(
          onPressed: () {
            showModalBottomSheet<bool>(
              context: context,
              backgroundColor: Colors.transparent,
              useRootNavigator: true,
              isScrollControlled: true,
              builder: (_) => _VotingOptionFormBottomSheet(
                title: 'Agregar opción de votación',
                actionLabel: 'Crear opción',
                onSubmit: (payload) =>
                    controller.createVotingOption(
                  votingId: voting.id,
                  data: payload,
                ),
              ),
            );
          },
          child: const Text('Agregar opción'),
        ),
      ],
    );
  }
}
