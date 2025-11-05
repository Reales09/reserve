import 'package:get/get.dart';
import 'package:rupu/domain/entities/horizontal_property_voting_groups.dart';
import 'package:rupu/domain/repositories/horizontal_properties_repository.dart';
import 'package:rupu/domain/infrastructure/repositories/horizontal_properties_repository_impl.dart';

class VotingGroupDetailController extends GetxController {
  final int propertyId;
  final int votingGroupId;
  final HorizontalPropertiesRepository repository;

  VotingGroupDetailController({
    required this.propertyId,
    required this.votingGroupId,
    HorizontalPropertiesRepository? repository,
  }) : repository = repository ?? HorizontalPropertiesRepositoryImpl();

  static String tagFor(int propertyId, int votingGroupId) =>
      'voting-group-detail-$propertyId-$votingGroupId';

  final group = Rxn<HorizontalPropertyVotingGroup>();
  final isLoading = false.obs;
  final errorMessage = RxnString();

  @override
  void onReady() {
    super.onReady();
    refresh();
  }

  Future<void> refresh() async {
    isLoading.value = true;
    errorMessage.value = null;
    try {
      final result = await repository.getVotingResults(
        propertyId: propertyId,
        groupId: votingGroupId,
        votingId: 0, // TODO: Get votingId from somewhere
      );
      if (result.success) {
        // TODO: Handle success
      } else {
        errorMessage.value = result.message ?? 'No se pudieron cargar los resultados de la votación.';
      }
    } catch (_) {
      errorMessage.value = 'No se pudo cargar la información del grupo de votación.';
    } finally {
      isLoading.value = false;
    }
  }
}
