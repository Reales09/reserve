import 'package:rupu/domain/entities/voting.dart';

class VotingsResult {
  final bool success;
  final String? message;
  final List<Voting> data;

  const VotingsResult({
    required this.success,
    this.message,
    required this.data,
  });
}
